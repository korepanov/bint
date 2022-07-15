float f1(){
	return 5.5;
};

string f2(string t){
	return "lalala";
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

int factorial(int n){
	int res;

	if (0 == n){
		return 1;	
	};
	
	res = (n * factorial((n - 1)));	
	
	return res;
};

void buy(string s){
	print(s);
	print("\n");
};

void main(){
	int res;
	hello();
	res = factorial(5);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
	buy("buy!");
};

main();
