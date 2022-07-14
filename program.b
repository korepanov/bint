goto(#factorial_end);
#factorial:
int n;
string $ret;
pop(n);
if (0 == n){
	push(1);
	$recurs_stack.pop($ret);
	goto($ret);	
};
push(n - 1);
$recurs_stack.push("#recurs_ret");
goto(#factorial);
#recurs_ret:
int n2;
pop(n2);
push(n * n2);
$recurs_stack.pop($ret);
goto($ret);

#factorial_end:
goto(#main_end);
#main:
int n; 
stack $recurs_stack;
push(5);
$recurs_stack.push("#here");
goto(#factorial);
#here:
pop(n)
string buf;
buf = str(n);
print(buf);
print("\n");
print("");
#main_end:
goto(#main);
