float descr(float a, float b, float c){
	return ((b^2) - ((4*a)*c));
};

stack solve(float a, float b, float c){
	float x1;
	float x2;
	stack res;
	float d; 
	d = descr(a, b, c);
	if (d < 0){
		return res;	
	};
	x1 = ((((-1)*b) - (d^0.5)) / (2 * a));
	x2 = ((((-1)*b) + (d^0.5)) / (2 * a));
	res.push(x2);
	res.push(x1);
	
	return res;
};
