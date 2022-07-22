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

int world(int num){
	if (5 == num){
		print("Hello world!\n");
		return num; 
	};
	
	num = (num + 1);
	world(num);
};

int buy(int num){
	if (5 == num){
		print("Hello world!\n");
		return num; 
	};
	
	num = (num + 1);
	buy(num);
};
void main(){
	int res;
	hello();
	res = factorial(5);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
	res = world(0);
	res = buy(0);
};

main();
