float compareStrings(string s, string s2){
	int l1;
	int l2; 

	l1 = len(s);
	l2 = len(s2);

	if (l2 > l1){
		string s3;
		s3 = s;
		s = s2;
		s2 = s3; 	
	};

	int sum;

	for (int i; i = 0; i < len(s2); i = (i + 1)){
		if (s[i] == s2[i]){
			sum = (sum + 1);			
		};					
	}; 

	return (float(sum) / float(len(s)));
	
};

void main(){
	print((str(compareStrings("something", "something")) + "\n"));
	
};
main();
