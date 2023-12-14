#import "stdlib/files.b"

void main(){
	int f;
	string s;	
	
	f = open_file("/home/slava/a1.txt", "read");	
	
	do{
		s = read_file(f, 10);
		print(s);	
	}while(NOT(s == ""));
	
	
	close_file(f);
	

};

main();
