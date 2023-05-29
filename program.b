void main(){
	print("Решение квадраного уравнения x^2+2x-3=0\n");
 	float a;
    a = 1;
	float b;
	b = 2;
	float c;
	c = (-3); 
	float d; 
	d = ((b^2.0)-(4*(a*c)));
	float x1;
	float x2; 

	x1 = ((((-1)*b)-(d^0.5))/(2*a));
    x2 = ((((-1)*b)+(d^0.5))/(2*a));

	print("x1 = ");
	print(str(x1));
	print("\n");
	print("x2 = "); 
	print(str(x2));
	print("\n");
};

main();
	



