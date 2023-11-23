#import "stdlib/core.b"

string root_source;
bool first_file;

void init(){
	first_file = True;
	root_source = "benv/prep_dowhile_program.b"; 
	SET_SOURCE(root_source);
	SET_DEST("benv/dowhile_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

void switch_files(){
	finish();	
	
	if (first_file){
		SET_SOURCE("benv/dowhile_program.b");
		SET_DEST("benv/dowhile_program2.b");
		first_file = False;			
	}else{
		SET_SOURCE("benv/dowhile_program2.b");
		SET_DEST("benv/dowhile_program.b");
		first_file = True;		
	};
};

void clear_files(){
	if (first_file){
		next_command(command);
		while(NOT(command == "end")){
			send_command(command);
			next_command(command);		
		};	
	};
		
	finish();
	DEL_DEST("benv/dowhile_program2.b");
	DEL_DEST(root_source); 
};

bool is_do(string command){
	stack s;
	string buf;

	buf = "do{";
	s = ops(command, buf);
	s.pop(buf);
	
	return (NOT("end" == buf));
};

string get_cond(string command){
	string op;
	stack s;
	string buf;
	int start_pos;
	int end_pos; 

	op = "while(";
	s = ops(command, op);
	s.pop(buf);
	start_pos = int(buf);
	start_pos = (start_pos + 5);
	end_pos = func_end(command, start_pos);
	start_pos = (start_pos + 1); 
	buf = command[start_pos:end_pos];

	return buf;
};

void main(){
	string buf;
	string cond; 
	int counter;
	bool was_here; 
	was_here = True;

	init();

	while (was_here){	
		was_here = False;
		switch_command();
		for (print(""); print(""); NOT(command == "end"); switch_command()){
			if (is_do(command)){
				was_here = True;
				counter = block_end(); 
				cond = get_cond(get_command(counter));

				buf = "for (print(\"\")";
				send_command(buf); 
				buf = "print(\"\")";
				send_command(buf);
				buf = "True";
				send_command(buf); 
				buf = "print(\"\")){print(\"\")";
				send_command(buf);
				
				switch_command();
				while (COMMAND_COUNTER < counter){
					send_command(command);					
					switch_command();		
				}; 

				buf = (("if (NOT(" + cond) + ")){print(\"\")");
				send_command(buf); 
				send_command("break");
				buf = "}"; 
				send_command(buf);
				send_command(buf); 
			}else{
				send_command(command); 
			};
		};
		switch_files(); 
		COMMAND_COUNTER = 0;
	};
	clear_files();  
	
};

main();
