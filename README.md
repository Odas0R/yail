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

The goal of the parser is to produce an AST.

# To solve operations +, -, / etc...

We need to use a top down operator precedence (or: PRATT PARSING).

> Looking at all these different forms of expressions it becomes clear that we
> need a really good approach to parse them correctly and in an understandable
> and extendable way. Our old 47approach of deciding what to do based on the
> current token won’t get us very far - at least not without wanting to tear
> our hair out. And that is where Vaughan Pratt comes in.

The parsing approach described by all three, which is called Top Down Operator Precedence
Parsing, or Pratt parsing, was invented as an alternative to parsers based on context-free gram-
mars and the Backus-Naur-Form.

**Terminology**:

- A `prefix` operator is an operator “in front of” its operand. Example: `--5`.
- A `postfix` operator is an operator “after” its operand. Example: `foobar++`.
- An `infix` operator sits between its operands, like this: `5 * 8`.

# We're writing a Pratt Parser

Pratt parser’s main idea is the association of parsing functions (which Pratt calls “semantic
code”) with token types. Whenever this token type is encountered, the parsing functions are
called to parse the appropriate expression and return an AST node that represents it. Each
token type can have up to two parsing functions associated with it, depending on whether the
token is found in a prefix or an infix position.

This is very neat when doing operations like +,<,>,/,- ...

But how does it work?

Well, in more simpler terms, it creates the AST Nodes via recursion, based on
the `token` types. It will get the current token type and execute the parser
for that expression.

Pratt doesn’t use a Parser structure and doesn’t pass around methods defined on `*Parser`. He
also doesn’t use maps and, of course, he didn’t use Go. His paper predates the release of Go by
36 years. And then there are naming differences: what we call prefixParseFns are “nuds” (for
“null denotations”) for Pratt. infixParseFns are “leds” (for “left denotations”).

**CHECK PAGE 69 for better understanding**
