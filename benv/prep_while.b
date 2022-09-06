#import "stdlib/core.b"

string root_source;

void init(){
	root_source = "benv/prep_for_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_while_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST(root_source);
};

bool is_while(string command){
	stack s;
	string buf;

	buf = "while(";
	s = ops(command, buf);
	s.pop(buf);
	
	return (NOT("end" == buf));
};

void main(){
	string buf;
	int pos;
	int command_len;
	stack s;

	init();
	#next:
	switch_command();
	if (NOT("end" == command)){
		if (is_while(command)){
			buf = "{";
			s = ops(command, buf);
			s.pop(buf);
			if ("end" == buf){
				println("prep_while: ERROR: no { after while");
				exit(1);			
			};
			pos = int(buf); 
			pos = (pos + 1);
			buf = command[0:pos];
			buf = (buf + "print(\"\")");
			send_command(buf);
			command_len = len(command);
			buf = command[pos:command_len];
			send_command(buf);		
		}else{
			send_command(command);
		};
		goto(#next);	
	};
	finish();
};

main();
