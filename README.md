A calculation framework for
1. Integers that couldn't be represented by built-in integer types (int, int32, int64 etc.)
2. Decimals (Since built-in float number types are not a precise representation)

The underlying data structure of an integer is an array of unsigned integers so that such an array could have a lot of bits to represent very large integers.
<div>The underlying data structure of a decimal is an integer generated by removing the decimal point of that decimal and an integer of type int to indicate where the point should be</div>

# Create an Integer
```go
b, err := CreateBinaryInteger("123")
```
When
1. The text is empty or
2. The text doesn't represent a decimal integer
<div>a non-nil error will be returned</div>
or an integer object will be created successfully

# Create a Decimal
```go
d, err := CreateNumber("123.45")
```
When the text doesn't represent a number, a non-nil error will be returned,
or a number object will be created successfully