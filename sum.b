int sum(int a, int b){
	return (a + b);
};

void sum2(int a, int b){
	int c;
	c = (a + b);
	string buf;
	buf = str(c);
	print(buf);
	print("\n");
};

int sum3(){
	return (2 + 3);
};

void sum4(){
	int a;
	a = 5;
	string buf;
	buf = str(a);
	print(buf);
	print("\n");
};

int res;
res = sum(2, 3);
string buf;
buf = str(res);
print(buf);
print("\n");
sum2(2, 3);
res = sum3();
buf = str(res);
print(buf);
print("\n");
sum4();
