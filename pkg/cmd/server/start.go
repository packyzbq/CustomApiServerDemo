package server

import (
	"CustomApiServerDemo/pkg/apiserver"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/apis/apiserver/v1alpha1"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/informers"
	"net"
)

/**
启动server步骤
1. 创建 config
2. complete config，填充非必填项
3. 创建 apiserver
4. 调用 server.Run
*/
const defaultEtcdPathPrefix = "/registry/restaurant.programming-kubernetes.info"

type CustomServerOptions struct {
	RecommendedOptions    *genericoptions.RecommendedOptions
	SharedInformerFactory informers.SharedInformerFactory
}

func NewCustomServerOptions(out, errOut io.Writer) *CustomServerOptions {
	o := &CustomServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			apiserver.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion),
			genericoptions.NewProcessInfo("pizza-apiserver", "pizza-apiserver"),
		),
	}

	return o
}

// NewCommandStartCustomServer provides a CLI handler for 'start master' command
// with a default CustomServerOptions.
func NewCommandStartCustomServer(defaults *CustomServerOptions, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch a custom API server",
		Long:  "Launch a custom API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)

	return cmd
}

func (o CustomServerOptions) Validate() error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *CustomServerOptions) Complete() error {
	// register admission plugins
	//pizzatoppings.Register(o.RecommendedOptions.Admission.Plugins)

	// add admisison plugins to the RecommendedPluginOrder
	//o.RecommendedOptions.Admission.RecommendedPluginOrder = append(o.RecommendedOptions.Admission.RecommendedPluginOrder, "PizzaToppings")

	return nil
}

func (o *CustomServerOptions) Config() (*apiserver.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	//o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
	//	client, err := clientset.NewForConfig(c.LoopbackClientConfig)
	//	if err != nil {
	//		return nil, err
	//	}
	//	informerFactory := informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
	//	o.SharedInformerFactory = informerFactory
	//	return []admission.PluginInitializer{custominitializer.New(informerFactory)}, nil
	//}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
	}
	return config, nil
}

func (o CustomServerOptions) Run(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	// 在apiserver启动之后，调用webhook启动 informer
	//_ = server.GenericAPIServer.AddPostStartHook("start-pizza-apiserver-informers", func(context genericapiserver.PostStartHookContext) error {
	//	config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
	//	o.SharedInformerFactory.Start(context.StopCh)
	//	return nil
	//})

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
