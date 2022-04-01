int hello(){

int a;
a = 0;

#start:
[print(""), (a < 5), goto(#end)];
a = (a + 1);
print("Hello world!\n");
goto(#start);
#end:
return 0;
};


int res;
res = hello();
[print(""), (0 == res), print("ERROR\n")];
