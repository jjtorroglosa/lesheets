# Tokens

```
TokenFrontMatter: "---" .* "---"
TokenHeader: "#".*\n
TokenHeaderBreak: "#-".*\n
TokenBar: "|" | "||:" | ":||"
TokenReturn: \n
TokenBarNote: "\"" .* "\""
TokenAnnotation: "!" .* "!"
TokenBacktick: "`" .* "`
TokenBacktickMultiline: "```" .* "```"
TokenEof: EOF
TokenChord: [^ ]+
```

# Parser
```
Song:
  FrontMatter Body
  |Body
  ;
FrontMatter: TokenFrontMatter
Body:
  Lines Sections
  |Sections
  ;
Sections:
  Section
  |Section Sections
Section:
  Header Lines
Header: TokenHeader
Lines:
  Line
  |Line Lines
  ;
Line:
  Bars TokenReturn
  |Bars TokenBar TokenReturn
  |TokenBar Bars TokenReturn
Bars:
  Bar
  |Bar TokenBar Bars
Bar:
  TokenBarNote TokenBar BarBody
  |TokenBar TokenBarNote BarBody
  |BarBody
  ;
BarBody:
  TokenBacktick
  |Chords
  ;
Chords:
  Chord
  |Chord Chords
  ;
Chord:
  TokenChord
  |TokenAnnotation TokenChord
  ;
```
