import HoganJsUtils from './hoganjs-utils';
import * as Rematch from './rematch';
import * as renderUtils from './render-utils';
import { DiffFile, DiffLine, DiffBlock, DiffLineDeleted, DiffLineContent, DiffLineContext, DiffLineInserted } from './types';
export interface LineByLineRendererConfig extends renderUtils.RenderConfig {
    renderNothingWhenEmpty?: boolean;
    matchingMaxComparisons?: number;
    maxLineSizeInBlockForComparison?: number;
}
export declare const defaultLineByLineRendererConfig: {
    renderNothingWhenEmpty: boolean;
    matchingMaxComparisons: number;
    maxLineSizeInBlockForComparison: number;
    matching: import("./types").LineMatchingType;
    matchWordsThreshold: number;
    maxLineLengthHighlight: number;
    diffStyle: import("./types").DiffStyleType;
};
export default class LineByLineRenderer {
    private readonly hoganUtils;
    private readonly config;
    constructor(hoganUtils: HoganJsUtils, config?: LineByLineRendererConfig);
    render(diffFiles: DiffFile[]): string;
    makeFileDiffHtml(file: DiffFile, diffs: string): string;
    generateEmptyDiff(): string;
    generateFileHtml(file: DiffFile): string;
    applyLineGroupping(block: DiffBlock): DiffLineGroups;
    applyRematchMatching(oldLines: DiffLine[], newLines: DiffLine[], matcher: Rematch.MatcherFn<DiffLine>): DiffLine[][][];
    processChangedLines(isCombined: boolean, oldLines: DiffLine[], newLines: DiffLine[]): FileHtml;
    generateLineHtml(oldLine?: DiffPreparedLine, newLine?: DiffPreparedLine): FileHtml;
    generateSingleLineHtml(line?: DiffPreparedLine): string;
}
declare type DiffLineGroups = [
    (DiffLineContext & DiffLineContent)[],
    (DiffLineDeleted & DiffLineContent)[],
    (DiffLineInserted & DiffLineContent)[]
][];
declare type DiffPreparedLine = {
    type: renderUtils.CSSLineClass;
    prefix: string;
    content: string;
    oldNumber?: number;
    newNumber?: number;
};
declare type FileHtml = {
    left: string;
    right: string;
};
export {};
//# sourceMappingURL=line-by-line-renderer.d.ts.map