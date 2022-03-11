int a;
bool basm;

int sum(int a, int b){
	int c; 
	c = (a + b);
	float myvar;
	bool my_bool_var;

	return c;
};
print("");
int sub(int a, int b){
	int c;
	c = (a - b);
	
	return c;
};
print("");
int mul(int a, int b){
	int c;
	c = (a * b);
	
	return c;
};

string buf;
a = sum(1, 2);
int b;
b = sub(a, 5);
b = mul(b, b);
buf = str(b); 
print(buf); 
print("\n");
