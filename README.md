# History

Logging of pages

[2023-03-14 00:36]: p.32
[2023-03-14 15:13]: p.38

> A parser is a software component that takes input data (frequently text) and
> builds a data structure -- often some kind of parse tree, abstract syntax
> tree or other hierarchical structure -- giving a structural representation of
> the input, checking for correct syntax in the process. [...] The Parser is
> often proceded by a separate lexical analyser, which creates tokens from the
> sequence of input characters;

A parser generator is a tool that, when fed with a formal description of a
language, produce parsers as their output. This output is the code that can be
compiled/interpreted and itself fed with source code as input to produce the
syntax tree.

Parsers differ in the format of the input they accept and the language of
output they produce. The majority of them use a _context-free grammar_ (CFG) as
their input. A CFG is a set of rules that describe how to form correctly (valid
according to the syntax) sentences in a language. The most common notational
formats are Backus-Naur Form (BNF) or Extended Backus-Naur Form (EBNF).

Example of BNF, EcmaScript:
```text
PrimaryExpression ::= "this"
                    | ObjectLiteral
                    | ( "(" Expression ")" )
                    | Identifier
                    | ArrayLiteral
                    | Literal

Literal ::= ( <DECIMAL_LITERAL>
            | <HEX_INTEGER_LITERAL>
            | <STRING_LITERAL>
            | <BOOLEAN_LITERAL>
            | <NULL_LITERAL>
            | <REGULAR_EXPRESSION_LITERAL> )

Identifier ::= <IDENTIFIER_NAME>

ArrayLiteral ::= "[" ( ( Elision )? "]"
                 | ElementList Elision "]"
                 | ( ElementList )? "]" )

ElementList ::= ( Elision )? AssignmentExpression
                ( Elision AssignmentExpression )*

Elision ::= ( "," )+

ObjectLiteral ::= "{" ( PropertyNameAndValueList )? "}"

PropertyNameAndValueList ::= PropertyNameAndValue ( "," PropertyNameAndValue
                                                  | "," )*

PropertyNameAndValue ::= PropertyName ":" AssignmentExpression

PropertyName ::= Identifier
              | <STRING_LITERAL>
              | <DECIMAL_LITERAL>
```

There are two main strategies:

1. top-down parsing
2. bottom-up parsing

the difference between top down and bottom up parsers is that the former starts
with constructing root node of the AST and then descends while the latter does
it the other way around.

A recursive descent parser, which works from the top down, is often recommended
for newcomers to parsing, since 31 it closely mirrors the way we think about
ASTs and their construction.

The parser produces an AST.




