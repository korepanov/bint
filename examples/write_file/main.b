// Программа записи в файл 

#import "stdlib/files.b"

void main(){
	int f;
	f = open_file("test.txt", "write");
	write_file(f, "Hello world!\n");
	close_file(f);
};

main();
