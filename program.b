float f1(){
	print("");
	return 5.5;
};

string f2(string t){
	print("");
	return "lalala";
};

int factorial(int n){
	print("");
	int res;

	if (0 == n){
		print("");
		return 1;	
	};
	
	res = (n * factorial((n - 1)));	
	
	return res;
};

stack f3(int y, int t){
	print("");
	stack res;
	res.push(1);
	return res;
};

bool f4(bool a, bool b, bool c){
	print("");
	return True; 
};


void hello(){
	print("");
	print("Hello!\n");
};

int world(int num, float y, int z){
	print("");
	if (5 == num){
		print("");
		print("Hello world!\n");
		return num; 
	};
	
	num = (num + 1);
	num = world(num, 5.5, 6);
};

int buy(int num, float y){
	print("");
	int t;
	
	if (5 == num){
		print("");
		print("Hello world!\n");
		return num; 
	};
	
	t = world(5, 7.7, 7);
	
	num = (num + 1);
	num = buy(num, 0);
	
	if (4 == num){
		print("");
		num = (num + 1);
		return buy(num, 0);
	};
};

int test(int num){
	print("");
	if (5 == num){
		print("");
		print("Hello world!\n");
		return num; 
	};
	
	num = (num + 1);
	num = test(num);
	
	if (4 == num){
		print("");
		num = (num + 1);
		return test(num);
	};
};

void inf(){
	print("inf");
	print("\n");
	inf();
};

void main(){
	print("");
	int res;
	hello();
	res = factorial(5);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
	res = world(1, 2, 3);
	res = buy(1, 2);
	inf();
};

main();
