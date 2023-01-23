#import "strings.b"
#import "io.b"
#import "solver.b"

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
	
	if ("end" == buf){
		println("No solution");	
	}else{
		float x1;
		float x2;
		x1 = float(buf);
		res.pop(x2);
		output(x1, x2); 
	};
	
};

main();
