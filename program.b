int sum(int a, int b){
	return (a + b);
};

int sum2(int a, int b){
	return (a + b);
};

void main(){
	int res;
	res = sum2(sum(1, 2), 3);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};

main();
