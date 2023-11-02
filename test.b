int i;
i = 0;
#for0:
if (i < 5){
	string s;
	print("Hello world!\n");
	i = (i + 1); 
};
if (i < 5){
	goto(#for0);
};
