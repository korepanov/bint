int sum(int a, int b){
	int c; 
	c = (a + b);
	return c; 
};

int sub(int a, int b){
	int c; 
	c = (a - b);
	return c; 
};

int a; 
a = sum(1, 2);
a = sub(a, 1);
string buf;
buf = str(a);
print(buf);
