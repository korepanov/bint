int sum(int a, int b){
	return (a + b);
};

int diff(int a, int b){
	return (a - b);
};

int mul(int a, int b){
	return (a * b);
};

float div(float a, float b){
	return (a / b);
};

void main(){
	float res;
	res = div(5, 2);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};

main();
