package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samber/lo"
)

const header = `// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

`

const headerTypescript = `// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

`

const bindingTemplate = `
/**Comments 
 * @function {{methodName}}* @param names {string}
 * @returns {Promise<string>}
 **/
`
const bindingTemplateTypescript = `Comments`

const callByIDTypescript = `export async function {{methodName}}({{inputs}}) : {{ReturnType}} {
	return Call.ByID({{ID}}{{params}});
}

`

const callByNameTypescript = `export async function {{methodName}}({{inputs}}) : {{ReturnType}} {
	return Call.ByName("{{Name}}"{{params}});
}

`

const callByID = `export async function {{methodName}}({{inputs}}) {
	return Call.ByID({{ID}}, ...Array.prototype.slice.call(arguments, 0));
}
`

const callByName = `export async function {{methodName}}({{inputs}}) {
	return Call.ByName("{{Name}}", ...Array.prototype.slice.call(arguments, 0));
}
`

var reservedWords = []string{
	"abstract",
	"arguments",
	"await",
	"boolean",
	"break",
	"byte",
	"case",
	"catch",
	"char",
	"class",
	"const",
	"continue",
	"debugger",
	"default",
	"delete",
	"do",
	"double",
	"else",
	"enum",
	"eval",
	"export",
	"extends",
	"false",
	"final",
	"finally",
	"float",
	"for",
	"function",
	"goto",
	"if",
	"implements",
	"import",
	"in",
	"instanceof",
	"int",
	"interface",
	"let",
	"long",
	"native",
	"new",
	"null",
	"package",
	"private",
	"protected",
	"public",
	"return",
	"short",
	"static",
	"super",
	"switch",
	"synchronized",
	"this",
	"throw",
	"throws",
	"transient",
	"true",
	"try",
	"typeof",
	"var",
	"void",
	"volatile",
	"while",
	"with",
	"yield",
	"object",
}

func sanitiseJSVarName(name string) string {
	// if the name is a reserved word, prefix with an
	// underscore
	if lo.Contains(reservedWords, name) {
		return "_" + name
	}
	return name
}

type ExternalStruct struct {
	Package string
	Name    string
}

func (p *Project) GenerateBinding(thisStructName string, method *BoundMethod, useIDs bool) (string, []string, map[packagePath]map[string]*ExternalStruct) {
	var externalStructs = make(map[packagePath]map[string]*ExternalStruct)
	var models []string
	template := bindingTemplate
	if useIDs {
		template += callByID
	} else {
		template += callByName
	}
	result := strings.ReplaceAll(template, "{{structName}}", thisStructName)
	result = strings.ReplaceAll(result, "{{methodName}}", method.Name)
	result = strings.ReplaceAll(result, "{{ID}}", fmt.Sprintf("%v", method.ID))

	// get last part of method.Package path
	parts := strings.Split(method.Package, "/")
	packageName := parts[len(parts)-1]

	result = strings.ReplaceAll(result, "{{Name}}", fmt.Sprintf("%v.%v.%v", packageName, thisStructName, method.Name))
	comments := strings.TrimSpace(method.DocComment)
	if comments != "" {
		comments = "\n * " + comments
	}
	result = strings.ReplaceAll(result, "Comments", comments)
	var params string
	for _, input := range method.JSInputs() {
		input.project = p
		inputName := sanitiseJSVarName(input.Name)
		pkgName := getPackageName(input)
		if pkgName != "" {
			models = append(models, pkgName)
		}
		if input.Type.IsStruct || input.Type.IsEnum {
			if _, ok := externalStructs[input.Type.Package]; !ok {
				externalStructs[input.Type.Package] = make(map[string]*ExternalStruct)
			}
			externalStructs[input.Type.Package][input.Type.Name] = &ExternalStruct{
				Package: input.Type.Package,
				Name:    input.Type.Name,
			}
		}

		inputType := input.JSType(packageName)
		params += "\n * @param " + inputName + " {" + inputType + "}"
	}
	params = strings.TrimSuffix(params, "\n")
	//if len(params) > 0 {
	//	params = "\n" + params
	//}
	result = strings.ReplaceAll(result, "* @param names {string}", params)
	var inputs string
	for _, input := range method.JSInputs() {
		pkgName := getPackageName(input)
		if pkgName != "" {
			models = append(models, pkgName)
		}
		inputs += sanitiseJSVarName(input.Name) + ", "
	}
	inputs = strings.TrimSuffix(inputs, ", ")
	args := inputs
	if len(args) > 0 {
		args = ", " + args
	}
	result = strings.ReplaceAll(result, "{{inputs}}", inputs)
	result = strings.ReplaceAll(result, "{{args}}", args)

	// outputs
	var returns string
	if len(method.Outputs) == 0 {
		returns = " * @returns {Promise<void>}"
	} else {
		returns = " * @returns {Promise<"
		for _, output := range method.Outputs {
			output.project = p
			pkgName := getPackageName(output)
			if pkgName != "" {
				models = append(models, pkgName)
			}
			jsType := output.JSType(pkgName)
			if jsType == "error" {
				jsType = "void"
			}
			if output.Type.IsStruct {
				if _, ok := externalStructs[output.Type.Package]; !ok {
					externalStructs[output.Type.Package] = make(map[string]*ExternalStruct)
				}
				externalStructs[output.Type.Package][output.Type.Name] = &ExternalStruct{
					Package: output.Type.Package,
					Name:    output.Type.Name,
				}
				jsType = output.NamespacedStructVariable(output.Type.Package)
			}
			returns += jsType + ", "
		}
		returns = strings.TrimSuffix(returns, ", ")
		returns += ">}"
	}
	result = strings.ReplaceAll(result, " * @returns {Promise<string>}", returns)

	return result, lo.Uniq(models), externalStructs
}

func (p *Project) GenerateBindingTypescript(thisStructName string, method *BoundMethod, useIDs bool) (string, []string, map[packagePath]map[string]*ExternalStruct) {
	var externalStructs = make(map[packagePath]map[string]*ExternalStruct)
	var models []string
	template := bindingTemplateTypescript
	if useIDs {
		template += callByIDTypescript
	} else {
		template += callByNameTypescript
	}
	result := strings.ReplaceAll(template, "{{structName}}", thisStructName)
	result = strings.ReplaceAll(result, "{{methodName}}", method.Name)
	result = strings.ReplaceAll(result, "{{ID}}", fmt.Sprintf("%v", method.ID))

	// get last part of method.Package path
	parts := strings.Split(method.Package, "/")
	packageName := parts[len(parts)-1]

	result = strings.ReplaceAll(result, "{{Name}}", fmt.Sprintf("%v.%v.%v", packageName, thisStructName, method.Name))
	comments := strings.TrimSpace(method.DocComment)
	if comments != "" {
		comments = "// " + comments + "\n"
	}
	result = strings.ReplaceAll(result, "Comments", comments)
	var params string
	for _, input := range method.JSInputs() {
		input.project = p
		inputName := sanitiseJSVarName(input.Name)
		pkgName := getPackageName(input)
		if pkgName != "" {
			models = append(models, pkgName)
		}
		if input.Type.IsStruct || input.Type.IsEnum {
			if _, ok := externalStructs[input.Type.Package]; !ok {
				externalStructs[input.Type.Package] = make(map[string]*ExternalStruct)
			}
			externalStructs[input.Type.Package][input.Type.Name] = &ExternalStruct{
				Package: input.Type.Package,
				Name:    input.Type.Name,
			}
		}
		params += ", " + inputName
	}
	result = strings.ReplaceAll(result, "{{params}}", params)
	//if len(params) > 0 {
	//	params = "\n" + params
	//}
	var inputs string
	for _, input := range method.JSInputs() {
		pkgName := getPackageName(input)
		if pkgName != "" {
			models = append(models, pkgName)
		}
		inputs += sanitiseJSVarName(input.Name) + ": " + input.JSType(packageName) + ", "
	}
	inputs = strings.TrimSuffix(inputs, ", ")
	args := inputs
	if len(args) > 0 {
		args = ", " + args
	}
	result = strings.ReplaceAll(result, "{{inputs}}", inputs)
	result = strings.ReplaceAll(result, "{{args}}", args)

	// outputs
	var returns string
	switch {
	case len(method.Outputs) == 0:
		returns = "Promise<void>"
	case len(method.Outputs) == 1 && method.Outputs[0].Type.Name == "error":
		returns = "Promise<void>"
	default:
		returns = "Promise<"
		for idx, output := range method.Outputs {
			output.project = p
			pkgName := getPackageName(output)
			if pkgName != "" {
				models = append(models, pkgName)
			}
			jsType := output.JSType(pkgName)
			if jsType == "error" {
				if len(method.Outputs) == 2 && idx == 1 {
					continue
				}

				jsType = "void"
			}
			if output.Type.IsStruct {
				if _, ok := externalStructs[output.Type.Package]; !ok {
					externalStructs[output.Type.Package] = make(map[string]*ExternalStruct)
				}
				externalStructs[output.Type.Package][output.Type.Name] = &ExternalStruct{
					Package: output.Type.Package,
					Name:    output.Type.Name,
				}
				jsType = output.NamespacedStructVariable(output.Type.Package)
			}
			returns += jsType + "|"
		}
		returns = strings.TrimSuffix(returns, "|")
		returns += ">"
	}
	result = strings.ReplaceAll(result, "{{ReturnType}}", returns)

	return result, lo.Uniq(models), externalStructs
}

func getPackageName(input *Parameter) string {
	if !input.Type.IsStruct {
		return ""
	}
	result := input.Type.Package
	if result == "" {
		result = "main"
	}
	return result
}

func isContext(input *Parameter) bool {
	return input.Type.Package == "context" && input.Type.Name == "Context"
}

func (p *Project) GenerateBindings(bindings map[string]map[string][]*BoundMethod, useIDs bool, useTypescript bool) map[string]map[string]string {

	var result = make(map[string]map[string]string)

	// sort the bindings keys
	packageNames := lo.Keys(bindings)
	sort.Strings(packageNames)
	for _, packageName := range packageNames {
		var allModels []string

		packageBindings := bindings[packageName]
		structNames := lo.Keys(packageBindings)
		relativePackageDir := p.RelativePackageDir(packageName)
		_ = relativePackageDir
		sort.Strings(structNames)
		for _, structName := range structNames {
			if _, ok := result[relativePackageDir]; !ok {
				result[relativePackageDir] = make(map[string]string)
			}
			methods := packageBindings[structName]
			sort.Slice(methods, func(i, j int) bool {
				return methods[i].Name < methods[j].Name
			})
			var allNamespacedStructs map[packagePath]map[string]*ExternalStruct
			var namespacedStructs map[packagePath]map[string]*ExternalStruct
			var thisBinding string
			var models []string
			var mainImports = ""
			if len(methods) > 0 {
				mainImports = "import {Call} from '@wailsio/runtime';\n"
			}
			for _, method := range methods {
				if useTypescript {
					thisBinding, models, namespacedStructs = p.GenerateBindingTypescript(structName, method, useIDs)
				} else {
					thisBinding, models, namespacedStructs = p.GenerateBinding(structName, method, useIDs)
				}
				// Merge the namespaced structs
				allNamespacedStructs = mergeNamespacedStructs(allNamespacedStructs, namespacedStructs)
				allModels = append(allModels, models...)
				result[relativePackageDir][structName] += thisBinding
			}

			if len(allNamespacedStructs) > 0 {
				thisPkg := p.packageCache[packageName]
				if !useTypescript {
					typedefs := "/**\n"
					for externalPackageName, namespacedStruct := range allNamespacedStructs {
						pkgInfo := p.packageCache[externalPackageName]
						relativePackageDir := p.RelativeBindingsDir(thisPkg, pkgInfo)
						namePrefix := ""
						if pkgInfo.Name != "" && pkgInfo.Path != thisPkg.Path {
							namePrefix = pkgInfo.Name
						}

						// Get keys from namespacedStruct and iterate over them in sorted order
						namespacedStructNames := lo.Keys(namespacedStruct)
						sort.Strings(namespacedStructNames)
						for _, thisStructName := range namespacedStructNames {
							structInfo := namespacedStruct[thisStructName]
							typedefs += " * @typedef {import('" + relativePackageDir + "/models')." + thisStructName + "} " + namePrefix + structInfo.Name + "\n"
						}
					}
					typedefs += " */\n"
					result[relativePackageDir][structName] = typedefs + result[relativePackageDir][structName]
				} else {
					// Generate imports instead of typedefs
					imports := ""
					for externalPackageName, namespacedStruct := range allNamespacedStructs {
						pkgInfo := p.packageCache[externalPackageName]
						relativePackageDir := p.RelativeBindingsDir(thisPkg, pkgInfo)
						namePrefix := ""
						if pkgInfo.Name != "" && pkgInfo.Path != thisPkg.Path {
							namePrefix = pkgInfo.Name
						}

						// Get keys from namespacedStruct and iterate over them in sorted order
						namespacedStructNames := lo.Keys(namespacedStruct)
						sort.Strings(namespacedStructNames)
						for _, thisStructName := range namespacedStructNames {
							structInfo := namespacedStruct[thisStructName]
							if namePrefix != "" {
								imports += "import {" + thisStructName + " as " + namePrefix + structInfo.Name + "} from '" + relativePackageDir + "/models';\n"
							} else {
								imports += "import {" + thisStructName + "} from '" + relativePackageDir + "/models';\n"
							}
						}
					}
					imports += "\n"
					result[relativePackageDir][structName] = imports + result[relativePackageDir][structName]
				}
			}
			if useTypescript {
				result[relativePackageDir][structName] = headerTypescript + mainImports + result[relativePackageDir][structName]
			} else {
				result[relativePackageDir][structName] = header + mainImports + result[relativePackageDir][structName]
			}
		}
	}

	return result
}

func mergeNamespacedStructs(structs map[packagePath]map[string]*ExternalStruct, structs2 map[packagePath]map[string]*ExternalStruct) map[packagePath]map[string]*ExternalStruct {
	if structs == nil {
		structs = make(map[packagePath]map[string]*ExternalStruct)
	}
	for pkg, pkgStructs := range structs2 {
		if _, ok := structs[pkg]; !ok {
			structs[pkg] = make(map[string]*ExternalStruct)
		}
		for name, structInfo := range pkgStructs {
			structs[pkg][name] = structInfo
		}
	}
	return structs
}
