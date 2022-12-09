int sum(int a, int b){
	return (a + b);
};

int diff(int a, int b){
	return (a - b);
};

int mul(int a, int b){
	return (a * b);
};

int div(int a, int b){
	return (a @ b);
};

void main(){
	int res;
	res = (diff(sum(3, 4), sum(1, 2)) + mul(diff(div(sum(1, 2), 3), 2), -5));
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};

main();
