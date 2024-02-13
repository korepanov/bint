// Программа чтения файла и вывода его на экран 
#import "stdlib/files.b"

void main(){
	int f; 
	string buf;
	string s;

	f = open_file("/home/slava/books/cookbook.txt", "read");	
		
	do{
		buf = read_file(f, 4096);
		s = (s + buf);	
	}while(NOT(buf==""));
	
	print(s);
	close_file(f);
	

};

main();
