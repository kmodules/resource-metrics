package cmds

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	resourcemetrics "kmodules.xyz/resource-metrics"
	"kmodules.xyz/resource-metrics/api"

	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
)

func NewCmdCalculate(clientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	var apiGroups []string
	cmd := &cobra.Command{
		Use:                   "calculate",
		Short:                 "Calculate metrics of a specific group of resources",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(clientGetter, sets.NewString(apiGroups...))
		},
	}
	cmd.Flags().StringSliceVar(&apiGroups, "apiGroups", apiGroups, "api groups for which to calculate resource")

	return cmd
}

func run(clientGetter genericclioptions.RESTClientGetter, apiGroups sets.String) error {
	cfg, err := clientGetter.ToRESTConfig()
	if err != nil {
		return err
	}
	client, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}
	mapper, err := clientGetter.ToRESTMapper()
	if err != nil {
		return err
	}

	r2 := map[schema.GroupVersionKind]core.ResourceList{}
	for _, gvk := range api.RegisteredTypes() {
		if apiGroups.Len() > 0 && !apiGroups.Has(gvk.Group) {
			continue
		}
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if meta.IsNoMatchError(err) {
			r2[gvk] = nil // keep track
			continue
		} else if err != nil {
			return err
		}

		var ri dynamic.ResourceInterface
		if mapping.Scope == meta.RESTScopeNamespace {
			ri = client.Resource(mapping.Resource).Namespace(core.NamespaceAll)
		} else {
			ri = client.Resource(mapping.Resource)
		}
		if result, err := ri.List(context.TODO(), metav1.ListOptions{}); err != nil {
			return err
		} else {
			var summary core.ResourceList
			for _, item := range result.Items {
				rr, err := resourcemetrics.AppResourceLimits(item.UnstructuredContent())
				if err != nil {
					return err
				}
				summary = api.AddResourceList(summary, rr)
			}
			r2[gvk] = summary
		}
	}

	gvks := make([]schema.GroupVersionKind, 0, len(r2))
	for gvk := range r2 {
		gvks = append(gvks, gvk)
	}
	sort.Slice(gvks, func(i, j int) bool {
		if gvks[i].Group == gvks[j].Group {
			return gvks[i].Kind < gvks[j].Kind
		}
		return gvks[i].Group < gvks[j].Group
	})

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.TabIndent)
	_, _ = fmt.Fprintln(w, "API VERSION\tKIND\tCPU\tMEMORY\tSTORAGE\t")
	for _, gvk := range gvks {
		rr := r2[gvk]
		if rr == nil {
			_, _ = fmt.Fprintf(w, "%s\t%s\t-\t-\t-\t\n", gvk.GroupVersion(), gvk.Kind)
		} else {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", gvk.GroupVersion(), gvk.Kind, rr.Cpu(), rr.Memory(), rr.Storage())
		}
	}
	return w.Flush()
}
