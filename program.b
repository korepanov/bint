float f1(){
	return 5.5;
};

string f2(string t){
	return "lalala";
};

int factorial(int n){
	int res;

	if (0 == n){
		return 1;	
	};
	
	res = (n * factorial((n - 1)));	
	
	return res;
};

stack f3(int y, int t){
	stack res;
	res.push(1);
	return res;
};

bool f4(bool a, bool b, bool c){
	return True; 
};


void hello(){
	print("Hello!\n");
};

int world(int num, float y, int z){
	if (5 == num){
		print("Hello world!\n");
		return num; 
	};

	string t;
	
	num = (num + 1);
	num = world(num, 5.5, 6);
	return num;
};

int buy(int num, float y){
	int t;
	
	if (5 == num){
		print("Hello world!\n");
		return num; 
	};
	
	t = world(5, 7.7, 7);
	num = (num + 1);
	string buf;
	buf = str(num);
	
	num = buy(num, 0);
	num = buy(num, 1);
	bool p;

	if (4 == num){
		num = (num + 1);
		return buy(num, 0);
	};
	
	return num;
};

int test(int num){
	if (5 == num){
		print("");
		print("Hello world!\n");
		return num; 
	};
	
	num = (num + 1);
	num = test(num);
	
	if (4 == num){
		num = (num + 1);
		return test(num);
	};
};

void inf(){
	print("inf");
	string p;
	print("\n");
	inf();
};

void main(){
	int res;
	hello();
	res = factorial(10);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};

main();
