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
	res = ((sum(1, sum(1, 2)) + mul(5, div(mul(sum(1, 2), 3), 4))) @ 5);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
	res = sum(((-1)*5), sum(sum(sum(1, 10), sum(2, 9)), 4));
	buf = str(res);
	print(buf);
	print("\n");
};

main();
