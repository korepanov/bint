#import "strings.b" 

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
