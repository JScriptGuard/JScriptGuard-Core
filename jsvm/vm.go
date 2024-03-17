package jsvm

import (
	"fmt"
	"log"

	"github.com/dop251/goja"
)

type JsVM struct {
	sharedFunctions map[string]SharedFunctionInterface
	cache           map[string]interface{}
	allowedBuckets  []string
	config          *JsVMConfig
	gojaVM          *goja.Runtime
	exports         *goja.Object
	scriptLoaded    bool
}

func (o *JsVM) prepareVM() error {
	vnh1Obj := o.gojaVM.NewObject()
	vnh1Obj.Set("com", o.gojaCOMFunctionModule)
	o.gojaVM.Set("vnh1", vnh1Obj)
	o.gojaVM.Set("exports", o.exports)
	return nil
}

func (o *JsVM) RunScript(script string) error {
	if o.scriptLoaded {
		return fmt.Errorf("LoadScript: always script loaded")
	}
	o.scriptLoaded = true
	_, err := o.gojaVM.RunString(script)
	if err != nil {
		panic(err)
	}
	return nil
}

func (o *JsVM) consoleLog(output string) {
	log.Println(output)
}

func (o *JsVM) consoleError(output string) {
	log.Println(output)
}

func (o *JsVM) gojaCOMFunctionModule(call goja.FunctionCall) goja.Value {
	// Es wird ermittelt um welchen vorgang es sich handelt
	if len(call.Arguments) < 1 {
		return o.gojaVM.ToValue("invalid")
	}

	// Die jeweilige Funktion wird ermittelt
	switch call.Arguments[0].String() {
	// Konsolen funktionen
	case "console":
		return console_base(o.gojaVM, call, o)
	// Share Functions
	case "root":
		return root_base(o.gojaVM, call, o)
	// S3 Funktionen
	case "s3":
		// Es wird geprüft ob die S3 Funktionen verfügbar sind
		if !o.config.EnableS3 {
			return goja.Undefined()
		}

		// Die S3 Funktionen werden bereitgestellt
		return sthreeb_base(o.gojaVM, call, o)
	// Die Cache Funktionen werden bereitgesllt
	case "cache":
		// Es wird geprüft ob Cache Funktion verfügbar sind
		if !o.config.EnableCache {
			return goja.Undefined()
		}

		// Die Cache Funktionen werden bereitgestellt
		return cache_base(o.gojaVM, call, o)
	// Es handelt sich um ein Unbekanntes Modul
	default:
		return goja.Undefined()
	}
}

func NewVM(config *JsVMConfig) (*JsVM, error) {
	// Die GoJA VM wird erstellt
	gojaVM := goja.New()

	// Das Basisobjekt wird erzeugt
	var vmObject *JsVM
	if config == nil {
		vmObject = &JsVM{config: &defaultConfig, gojaVM: gojaVM, scriptLoaded: false, exports: gojaVM.NewObject(), sharedFunctions: make(map[string]SharedFunctionInterface), allowedBuckets: make([]string, 0), cache: make(map[string]interface{})}
	} else {
		vmObject = &JsVM{config: config, gojaVM: gojaVM, scriptLoaded: false, exports: gojaVM.NewObject(), sharedFunctions: make(map[string]SharedFunctionInterface), allowedBuckets: make([]string, 0), cache: make(map[string]interface{})}
	}

	// Die Funktionen werden hinzugefügt
	if err := vmObject.prepareVM(); err != nil {
		return nil, fmt.Errorf("NewVM: " + err.Error())
	}

	// Das VM Objekt wird zurückgegeben
	return vmObject, nil
}