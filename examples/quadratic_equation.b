void println(string s){
	print(s);
	print("\n");
};

float descr(float a, float b, float c){
	return ((b^2) - ((4*a)*c));
};

stack solve(float a, float b, float c){
	float x1;
	float x2;
	stack res;
	float d; 
	d = descr(a, b, c);
	[goto(#solve_e), (d < 0), print("")];
	x1 = ((((-1)*b) - (d^0.5)) / (2 * a));
	x2 = ((((-1)*b) + (d^0.5)) / (2 * a));
	res.push(x2);
	res.push(x1);
	#solve_e:
	return res;
};

stack input_data(){
	float a;
	float b;
	float c;
	string buf;
	
	println("Input a");
	input(buf);
	a = float(buf);
	
	println("Input b");
	input(buf);
	b = float(buf);
	
	println("Input c");
	input(buf);
	c = float(buf);
	
	stack res;
	res.push(c);
	res.push(b);
	res.push(a);
	return res;
		
};

void output(float x1, float x2){
	string buf;
	println("Solution:");
	print("x1=");
	buf = str(x1);
	println(buf);
	print("x2=");
	buf = str(x2);
	println(buf);
};
void main(){
	println("Solving equation of the form ax^2 + bx + c = 0");
	stack res;
	res = input_data();
	float a;
	float b;
	float c; 
	res.pop(a);
	res.pop(b);
	res.pop(c);
	res = solve(a, b, c);
	string buf;
	res.pop(buf);
	[goto(#no_solution), ("end" == buf), print("")];
	res.push(buf);
	float x1;
	float x2;
	res.pop(x1);
	res.pop(x2);
	output(x1, x2);
	goto(#main_e);
	#no_solution:
	println("No solution");
	#main_e:
	print("");
};

main();
