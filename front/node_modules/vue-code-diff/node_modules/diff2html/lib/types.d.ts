export declare type DiffLineParts = {
    prefix: string;
    content: string;
};
export declare enum LineType {
    INSERT = "insert",
    DELETE = "delete",
    CONTEXT = "context"
}
export interface DiffLineDeleted {
    type: LineType.DELETE;
    oldNumber: number;
    newNumber: undefined;
}
export interface DiffLineInserted {
    type: LineType.INSERT;
    oldNumber: undefined;
    newNumber: number;
}
export interface DiffLineContext {
    type: LineType.CONTEXT;
    oldNumber: number;
    newNumber: number;
}
export declare type DiffLineContent = {
    content: string;
};
export declare type DiffLine = (DiffLineDeleted | DiffLineInserted | DiffLineContext) & DiffLineContent;
export interface DiffBlock {
    oldStartLine: number;
    oldStartLine2?: number;
    newStartLine: number;
    header: string;
    lines: DiffLine[];
}
export interface DiffFileName {
    oldName: string;
    newName: string;
}
export interface DiffFile extends DiffFileName {
    addedLines: number;
    deletedLines: number;
    isCombined: boolean;
    isGitDiff: boolean;
    language: string;
    blocks: DiffBlock[];
    oldMode?: string | string[];
    newMode?: string;
    deletedFileMode?: string;
    newFileMode?: string;
    isDeleted?: boolean;
    isNew?: boolean;
    isCopy?: boolean;
    isRename?: boolean;
    isBinary?: boolean;
    isTooBig?: boolean;
    unchangedPercentage?: number;
    changedPercentage?: number;
    checksumBefore?: string | string[];
    checksumAfter?: string;
    mode?: string;
}
export declare type OutputFormatType = 'line-by-line' | 'side-by-side';
export declare const OutputFormatType: {
    [_: string]: OutputFormatType;
};
export declare type LineMatchingType = 'lines' | 'words' | 'none';
export declare const LineMatchingType: {
    [_: string]: LineMatchingType;
};
export declare type DiffStyleType = 'word' | 'char';
export declare const DiffStyleType: {
    [_: string]: DiffStyleType;
};
//# sourceMappingURL=types.d.ts.map