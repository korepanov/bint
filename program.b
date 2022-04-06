int summa(int a, int b){
	return (a + b);
};

int res;

int summa2(){
	return (2 + 3);
};

float subtitution(float first, float second){
	return (first - second);
};

res = (((100 + (((4 + summa(1, 2)) + summa(2, 3)) + summa(4, 5))) + 200) + summa(300, 400));
string buf;
buf = str(res);
print(buf);
print("\n");
