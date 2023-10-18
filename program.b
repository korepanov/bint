string cut(string s){
	int i;
	i = index(s, "Foundation"); 
	print((str(i) + "\n"));
	print((str(len(s)) + "\n"));
	return s[i:len(s)];
};

void main(){
	string s;
	s = "English Wikipedia is hosted alongside other language editions by the Wikimedia Foundation";
	s = (s + ", an American nonprofit organization. Its content is written independently of other");
	s = (s + " editions[1] in various varieties of English, aiming to stay consistent within");
	s = (s + " articles. Its internal newspaper is The Signpost."); 

	print((cut(s) + "\n"));
	 
	
};
main();
