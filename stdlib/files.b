/* opens file with file_path and mode 
 
 modes:
 mode = read 
 mode = write
 mode = append 

 sets error if failed
 return file descriptor number  
*/

int open_file(string file_path, string mode){
	int res; 
	int m;

	if (mode == "read"){
		m = 0;
	}else if (mode == "write"){
		m = 1;	
	}else if (mode = "append"){
		m = 2;	
	}else{
		error = "no such mode for opening file";
		return -1;
	};

	res = $open_f(file_path, m); 
	return res;
};

/*
 reads size bytes from file with descriptor_number into variable s

 sets error if failed 
 returns number of read bytes 
*/
int read_file(int descriptor_number, int size; string s){
	int bytes; 
	bytes = $read_f(descriptor_number, size, s);
	return bytes;
};

/*
 closes file with descriptor_number 

 sets error if failed
*/
void close_file(int descriptor_number){
	$close_f(f); 
	try{
		string s; 
		read_file(f, 1, s);	
	};
	if (NOT("" == error)){
		error = ("could not close file with file descriptor number " + str(descriptor_number));	
		
		if ($toPanic){
			print((error + "\n"));
			exit(1);		
		};	
	};
};
