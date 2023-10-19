string cut(string s){
	return (s[index(s, "Foundation"):len(s)]);
};

int pow(int a, int b){
	float c;
	c = float(a);
	float d;
	d = float(b);
	return (int((c^d)));
};

void main(){
	string s;
	s = "English Wikipedia is hosted alongside other language editions by the Wikimedia Foundation";
	s = (s + ", an American nonprofit organization. Its content is written independently of other");
	s = (s + " editions[1] in various varieties of English, aiming to stay consistent within");
	s = (s + " articles. Its internal newspaper is The Signpost."); 

	print((cut(s) + "\n"));
	print((str(pow(2, 4)) + "\n")); 	 
	
};
main();
