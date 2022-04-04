int solve(int a, int b){
	return ((a+b)*4);
};

void solve2(int a, int b){
	int c;
	c = ((a+b)*4);
	string buf;
	buf = str(c);
	print(buf);
	print("\n");
};

int solve3(){
	int a;
	int b;
	a = 2;
	b = 3;
	return ((a+b)*4);
};

void solve4(){
	string buf;
	int c; 
	c = ((a+b)*4);

	buf = str(c);
	print(buf);
	print("\n");
};

int res; 
res = solve(2, 3);
string buf;
buf = str(res);
print(buf);
print("\n");
solve2(2, 3);
res = solve3();
buf = str(res);
print(buf);
print("\n");
solve4(); 
