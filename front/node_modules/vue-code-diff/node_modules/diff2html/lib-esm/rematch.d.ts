export declare type BestMatch = {
    indexA: number;
    indexB: number;
    score: number;
};
export declare function levenshtein(a: string, b: string): number;
export declare type DistanceFn<T> = (x: T, y: T) => number;
export declare function newDistanceFn<T>(str: (value: T) => string): DistanceFn<T>;
export declare type MatcherFn<T> = (a: T[], b: T[], level?: number, cache?: Map<string, number>) => T[][][];
export declare function newMatcherFn<T>(distance: (x: T, y: T) => number): MatcherFn<T>;
//# sourceMappingURL=rematch.d.ts.map