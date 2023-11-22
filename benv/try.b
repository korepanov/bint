#import "stdlib/core.b"

string root_source;
int num;

void init(){
	root_source = "benv/prep_try_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/try_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST(root_source);
}; 

bool is_try(string command){ 
	stack st;
	string op; 
	op = "try{";
 	
	st = ops(command, op);
	string buf; 
	st.pop(buf); 
	if (buf == "end"){
		return False;	
	};

	int pos;
	pos = int(buf);
	
	if (NOT(pos == 0)){
		print("try.b: error in the try format\n");
		exit(1); 	
	};
	
	return True; 
};

void modify_block(){
	string snum;
	snum = str(num);		
	int e;
	e = block_end();
	
	string buf;
	buf = "error=\"\"";
	send_command(buf);

	buf = "$toPanic=False";
	send_command(buf); 
	
	
	buf = "if(True){print(\"\")";
	send_command(buf); 

	switch_command();
	while (COMMAND_COUNTER < e){
		buf = "if(NOT(error==\"\")){print(\"\")";
		send_command(buf);
		buf = (("goto(#_attempt" + snum) + ")");
		send_command(buf);
		buf = "}";
		send_command(buf);
		send_command(command);
		switch_command();			
	};

	
	buf = (("#_attempt" + snum) + ":\nprint(\"\")"); 
	send_command(buf);
	num = (num + 1);

	buf = "$toPanic=True";
	send_command(buf);  

};

void main(){
	init();
	
	switch_command();

	while (NOT(command == "end")){
		if (is_try(command)){
			modify_block();		
		}; 
		send_command(command); 
		switch_command(); 	
	};
	


	finish(); 
};

main();
