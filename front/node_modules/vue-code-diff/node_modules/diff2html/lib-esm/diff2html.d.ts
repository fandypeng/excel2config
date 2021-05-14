import * as DiffParser from './diff-parser';
import { LineByLineRendererConfig } from './line-by-line-renderer';
import { SideBySideRendererConfig } from './side-by-side-renderer';
import { DiffFile, OutputFormatType } from './types';
import { HoganJsUtilsConfig } from './hoganjs-utils';
export interface Diff2HtmlConfig extends DiffParser.DiffParserConfig, LineByLineRendererConfig, SideBySideRendererConfig, HoganJsUtilsConfig {
    outputFormat?: OutputFormatType;
    drawFileList?: boolean;
}
export declare const defaultDiff2HtmlConfig: {
    outputFormat: OutputFormatType;
    drawFileList: boolean;
    renderNothingWhenEmpty: boolean;
    matchingMaxComparisons: number;
    maxLineSizeInBlockForComparison: number;
    matching: import("./types").LineMatchingType;
    matchWordsThreshold: number;
    maxLineLengthHighlight: number;
    diffStyle: import("./types").DiffStyleType;
};
export declare function parse(diffInput: string, configuration?: Diff2HtmlConfig): DiffFile[];
export declare function html(diffInput: string | DiffFile[], configuration?: Diff2HtmlConfig): string;
//# sourceMappingURL=diff2html.d.ts.map