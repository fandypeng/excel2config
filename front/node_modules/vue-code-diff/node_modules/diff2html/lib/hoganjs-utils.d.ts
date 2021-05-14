import * as Hogan from 'hogan.js';
export interface RawTemplates {
    [name: string]: string;
}
export interface CompiledTemplates {
    [name: string]: Hogan.Template;
}
export interface HoganJsUtilsConfig {
    compiledTemplates?: CompiledTemplates;
    rawTemplates?: RawTemplates;
}
export default class HoganJsUtils {
    private preCompiledTemplates;
    constructor({ compiledTemplates, rawTemplates }: HoganJsUtilsConfig);
    static compile(templateString: string): Hogan.Template;
    render(namespace: string, view: string, params: Hogan.Context, partials?: Hogan.Partials, indent?: string): string;
    template(namespace: string, view: string): Hogan.Template;
    private templateKey;
}
//# sourceMappingURL=hoganjs-utils.d.ts.map