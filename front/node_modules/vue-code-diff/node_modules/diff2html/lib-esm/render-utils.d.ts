import { LineMatchingType, DiffStyleType, LineType, DiffLineParts, DiffFile, DiffFileName } from './types';
export declare type CSSLineClass = 'd2h-ins' | 'd2h-del' | 'd2h-cntx' | 'd2h-info' | 'd2h-ins d2h-change' | 'd2h-del d2h-change';
export declare const CSSLineClass: {
    [_: string]: CSSLineClass;
};
export declare type HighlightedLines = {
    oldLine: {
        prefix: string;
        content: string;
    };
    newLine: {
        prefix: string;
        content: string;
    };
};
export interface RenderConfig {
    matching?: LineMatchingType;
    matchWordsThreshold?: number;
    maxLineLengthHighlight?: number;
    diffStyle?: DiffStyleType;
}
export declare const defaultRenderConfig: {
    matching: LineMatchingType;
    matchWordsThreshold: number;
    maxLineLengthHighlight: number;
    diffStyle: DiffStyleType;
};
export declare function toCSSClass(lineType: LineType): CSSLineClass;
export declare function escapeForHtml(str: string): string;
export declare function deconstructLine(line: string, isCombined: boolean, escape?: boolean): DiffLineParts;
export declare function filenameDiff(file: DiffFileName): string;
export declare function getHtmlId(file: DiffFileName): string;
export declare function getFileIcon(file: DiffFile): string;
export declare function diffHighlight(diffLine1: string, diffLine2: string, isCombined: boolean, config?: RenderConfig): HighlightedLines;
//# sourceMappingURL=render-utils.d.ts.map