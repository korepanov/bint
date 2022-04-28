string$303root_source;
string$303first_command;
string$303$init_return_var;
goto$303 #init_end$303 null;
#init:
print$303 ""$303 null;
print$303 ""$303 null;
get_root_source$303 root_source$303 null;
SET_SOURCE$303 root_source$303 null;
SET_DEST$303 "bendBenv/import_program.b"$303 null;
goto$303 $init_return_var$303 null;
#init_end:
print$303 ""$303 null;
string$303$finish_return_var;
goto$303 #finish_end$303 null;
#finish:
print$303 ""$303 null;
print$303 ""$303 null;
UNSET_SOURCE;
UNSET_DEST;
goto$303 $finish_return_var$303 null;
#finish_end:
print$303 ""$303 null;
string$303$find_import_end_return_var;
int$303$find_import_end_res;
goto$303 #find_import_end_end$303 null;
#find_import_end:
print$303 ""$303 null;
string$303command;
print$303 ""$303 null;
pop$303 command$303 null;
print$303 ""$303 null;
int$303command_len;
int$303res;
=$303 command_len$303 len$303 command$303 null;
command=command[8:command_len];
=$303 res$303 index$303 command$303 """;
=$303 res$303 +$303 res$303 8;
push$303 res$303 null;
UNDEFINE$303 res$303 null;
UNDEFINE$303 command_len$303 null;
UNDEFINE$303 command$303 null;
goto$303 $find_import_end_return_var$303 null;
#find_import_end_end:
print$303 ""$303 null;
string$303$reverse_return_var;
stack$303$reverse_res;
goto$303 #reverse_end$303 null;
#reverse:
print$303 ""$303 null;
stack$303s;
print$303 ""$303 null;
pop$303 s$303 null;
print$303 ""$303 null;
string$303buf;
stack$303res;
s.pop(buf);
#reverse_s:
print$303 ""$303 null;
[goto(#reverse_e),("end"==buf),print("")];
res.push(buf);
s.pop(buf);
print$303 ""$303 null;
goto$303 #reverse_s$303 null;
#reverse_e:
print$303 ""$303 null;
push$303 res$303 null;
UNDEFINE$303 res$303 null;
UNDEFINE$303 buf$303 null;
UNDEFINE$303 s$303 null;
goto$303 $reverse_return_var$303 null;
#reverse_end:
print$303 ""$303 null;
string$303$get_imports_return_var;
stack$303$get_imports_res;
goto$303 #get_imports_end$303 null;
#get_imports:
print$303 ""$303 null;
print$303 ""$303 null;
stack$303res;
int$303number;
string$303import_name;
string$303command;
int$303command_len;
int$303import_end;
next_command$303 command$303 null;
#get_imports_s:
print$303 ""$303 null;
[goto(#get_imports_e),("end"==command),print("")];
print$303 ""$303 null;
=$303 number$303 index$303 command$303 "#import";
print$303 ""$303 null;
[print(""),(0==number),goto(#get_imports_e)];
=$303 command_len$303 len$303 command$303 null;
int$303$find_import_end_res0;
push$303 command$303 null;
=$303 $find_import_end_return_var$303 "#find_import_end_res0";
goto$303 #find_import_end$303 null;
#find_import_end_res0:
pop$303 $find_import_end_res$303 null;
=$303 $find_import_end_res0$303 $find_import_end_res;
=$303 import_end$303 $find_import_end_res0;
import_name=command[8:import_end];
res.push(import_name);
=$303 import_end$303 +$303 import_end$303 1;
command=command[import_end:command_len];
print$303 ""$303 null;
goto$303 #get_imports_s$303 null;
#get_imports_e:
print$303 ""$303 null;
RESET_SOURCE;
=$303 first_command$303 command;
stack$303$reverse_res0;
push$303 res$303 null;
=$303 $reverse_return_var$303 "#reverse_res0";
goto$303 #reverse$303 null;
#reverse_res0:
pop$303 $reverse_res$303 null;
=$303 $reverse_res0$303 $reverse_res;
push$303 $reverse_res0$303 null;
UNDEFINE$303 $reverse_res0$303 null;
UNDEFINE$303 $find_import_end_res0$303 null;
UNDEFINE$303 import_end$303 null;
UNDEFINE$303 command_len$303 null;
UNDEFINE$303 command$303 null;
UNDEFINE$303 import_name$303 null;
UNDEFINE$303 number$303 null;
UNDEFINE$303 res$303 null;
goto$303 $get_imports_return_var$303 null;
#get_imports_end:
print$303 ""$303 null;
string$303$file_union_return_var;
goto$303 #file_union_end$303 null;
#file_union:
print$303 ""$303 null;
print$303 ""$303 null;
string$303command;
stack$303imports;
string$303import;
stack$303$get_imports_res0;
=$303 $get_imports_return_var$303 "#get_imports_res0";
goto$303 #get_imports$303 null;
#get_imports_res0:
pop$303 $get_imports_res$303 null;
=$303 $get_imports_res0$303 $get_imports_res;
=$303 imports$303 $get_imports_res0;
imports.pop(import);
#file_union_s:
print$303 ""$303 null;
[goto(#file_union_e),("end"==import),print("")];
UNSET_SOURCE;
SET_SOURCE$303 import$303 null;
next_command$303 command$303 null;
#import_s:
print$303 ""$303 null;
[goto(#import_e),("end"==command),print("")];
send_command$303 command$303 null;
next_command$303 command$303 null;
print$303 ""$303 null;
goto$303 #import_s$303 null;
#import_e:
print$303 ""$303 null;
imports.pop(import);
print$303 ""$303 null;
goto$303 #file_union_s$303 null;
#file_union_e:
print$303 ""$303 null;
SET_SOURCE$303 root_source$303 null;
next_command$303 command$303 null;
send_command$303 first_command$303 null;
next_command$303 command$303 null;
#final_s:
print$303 ""$303 null;
[goto(#final_e),("end"==command),print("")];
send_command$303 command$303 null;
next_command$303 command$303 null;
print$303 ""$303 null;
goto$303 #final_s$303 null;
#final_e:
print$303 ""$303 null;
print$303 ""$303 null;
UNDEFINE$303 $get_imports_res0$303 null;
UNDEFINE$303 import$303 null;
UNDEFINE$303 imports$303 null;
UNDEFINE$303 command$303 null;
goto$303 $file_union_return_var$303 null;
#file_union_end:
print$303 ""$303 null;
string$303$main_return_var;
goto$303 #main_end$303 null;
#main:
print$303 ""$303 null;
print$303 ""$303 null;
=$303 $init_return_var$303 "#init_res0";
goto$303 #init$303 null;
#init_res0:
print$303 ""$303 null;
print$303 ""$303 null;
=$303 $file_union_return_var$303 "#file_union_res0";
goto$303 #file_union$303 null;
#file_union_res0:
print$303 ""$303 null;
print$303 ""$303 null;
=$303 $finish_return_var$303 "#finish_res0";
goto$303 #finish$303 null;
#finish_res0:
print$303 ""$303 null;
print$303 ""$303 null;
goto$303 $main_return_var$303 null;
#main_end:
print$303 ""$303 null;
=$303 $main_return_var$303 "#main_res0";
goto$303 #main$303 null;
#main_res0:
print$303 ""$303 null;
print$303 ""$303 null;
