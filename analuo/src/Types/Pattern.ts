enum PatternType {
  Unit = "unit",
  Range = "range",
  List = "list",
  NestedRange = "nested range",
}

export type Pattern = {
  expression: string;
  type: PatternType;
  map?: {
    delimiter?: number;
    descriptor?: number;
    startValue?: number;
    endValue?: number;
    unitValue?: number;
  };
};

export const basePatterns: Pattern[] = [
  {
    expression: "^[^a-z0-9]?(\\s*\\(.+\\))*$",
    type: PatternType.Unit,
  },
  {
    expression: "^([A-Z]|[0-9]+|[A-Z][0-9]+|[0-9]+[A-Z])(?:\\s*\\(?:.+\\))*$",
    type: PatternType.Unit,
    map: {
      unitValue: 2,
    },
  },
  {
    expression: "^(APT|APTS|UNIT|LOT|SUITE|BUILDING|COTTAGE #|TRAILER|GARAGE BUILDING)\\s*([A-Z0-9]+-*[A-Z0-9]*)(?:\\s*\\(?:.+\\))*$",
    type: PatternType.Unit,
    map: {
      descriptor: 2,
      unitValue: 3,
    },
  },
  {
    expression: "^(?:.*)(APT|APTS|UNIT|LOT)*\\s*([A-Z][0-9]*)\\s*-\\s*(?:APT|APTS|UNIT)*\\s*([A-Z][0-9]*)(\\s*\\(.+\\))*$",
    type: PatternType.Range,
    map: {
      descriptor: 2,
      startValue: 3,
      endValue: 4,
    },
  },
  {
    expression: "^(?:.*)(APT|APTS|UNIT|LOT|SUITE)\\s*([A-Z0-9]+)\\s*-\\s*(?:APT|APTS|UNIT|LOT|SUITE)\\s*([A-Z0-9]+)\\s*(?:\\s*\\(?:.+\\))*$",
    type: PatternType.Range,
    map: {
      descriptor: 2,
      startValue: 3,
      endValue: 4,
    },
  },
  {
    expression: "^(?:.+(,|;| AND ))+.+(?:\\s*\\(?:.+\\))*$",
    type: PatternType.List,
    map: {
      delimiter: 2,
    },
  },
  {
    expression: "EACH.* HAS AN",
    type: PatternType.NestedRange,
  },
  {
    expression: "(\\s|\\()OR ",
    type: PatternType.NestedRange,
  },
];

// export type Pattern = {
//   expression: string;
//   type: string;
//   map?: {
//     delimiter: number | undefined;
//     unitValue?: number | undefined;
//     descriptor?: number | undefined;
//     startValue?: number | undefined;
//     endValue?: number | undefined;
//   };
// }