string, $println_return_var;
goto, #println_end, null;
#println:
print, "", null;
string, s;
print, "", null;
pop, s, null;
print, "", null;
print, s, null;
print, "\n", null;
UNDEFINE, s, null;
goto, $println_return_var, null;
#println_end:
print, "", null;
string, $hello_user_return_var;
goto, #hello_user_end, null;
#hello_user:
print, "", null;
string, s;
print, "", null;
pop, s, null;
print, "", null;
=, s, +, +, "Hello, ", s, "!";
=, s, +, +, "Hello, ", s, "!";
=, s, +, +, "Hello, ", s, "!";
push, s, null;
=, $println_return_var, "#println_res0";
goto, #println, null;
#println_res0:
print, "", null;
print, "", null;
UNDEFINE, s, null;
goto, $hello_user_return_var, null;
#hello_user_end:
print, "", null;
string, $main_return_var;
goto, #main_end, null;
#main:
print, "", null;
print, "", null;
string, s;
input, s, null;
push, s, null;
=, $hello_user_return_var, "#hello_user_res0";
goto, #hello_user, null;
#hello_user_res0:
print, "", null;
print, "", null;
UNDEFINE, s, null;
goto, $main_return_var, null;
#main_end:
print, "", null;
=, $main_return_var, "#main_res0";
goto, #main, null;
#main_res0:
print, "", null;
print, "", null;
