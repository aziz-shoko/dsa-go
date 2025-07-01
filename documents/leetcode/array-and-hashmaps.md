## Decision Tree: Array Vs HashMaps vs Both
```
Is the problem asking for:
├── Position/Index based operations? → Array manipulation
├── Frequency/Counting? → HashMap
├── Fast lookups? → HashMap
├── Grouping by property? → HashMap
├── Finding pairs/complements? → HashMap + Array iteration
├── Maintaining order while checking duplicates? → Both
└── Complex pattern matching? → Usually HashMap + Array
```

## Problem Solving Methodology

### Step 1: Pattern Recognition Question
* Do I need to remember what I've seen before?
* Am I looking for dpulicates or freqeuncies?
* Do I need fast lookups?
* Am I grouping things by a property

### Step 2: Data Structure Selection
* Arrays only: Direct indexing, simple iterations
* HashMap only: Pure counting, grouping, existence checks
* Both: Position matters and need fast lookups

### Step 3: Algorithm Design
1. Identify what to store (element, index, count, computed keys)
2. Determine when to store (before or after check)
3. Define the lookup condition
4. Handle edge cases (empty input, single element)


## Time & Space Complexity Patterns

HashMap Operations:
* Access/Insert/Delete: O(1) average, O(n) worse case
* Space: O(k) where k is number of unique elements

Common Complexity Profiles
* Single pass with HashMap: O(n) time, O(k) space
* Two pass: O(n) time, O(k) space
* Nested lookups: Watch for O(n^2) time complexity