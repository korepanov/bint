bool func1(){
	int a;
	int b;
	b = 1;
	string r;
	r = "NOT";
	bool p;
	p = ("NOT" == r);
	return (NOT((NOT(NOT(((a < b)XOR(b < a))AND(False))))XOR(NOT(4 < 5))));	
};

void main(){
	bool t;
	t = func1();
	string s;
	s = str(t);
	print(s); 
	
};

main();
