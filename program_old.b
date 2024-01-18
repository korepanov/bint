#import "stdlib/files.b"

void main(){
	int f;
	string s;	
	string buf; 
	
	f = open_file("/home/slava/books/cookbook.txt", "read");	
	
	do{
		buf = read_file(f, 1);
		print(buf);	
	}while(True);
	
	print(s);
	close_file(f);
	

};

main();
