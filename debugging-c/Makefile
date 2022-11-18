
NAME      		= bugged-program

# CC          	= clang++

ifdef DEBUG
CFLAGS      	= -Wall -Wextra -Werror -O0 -g -fsanitize=address
else
CFLAGS 			= -Wall -Wextra -Werror -O3
endif

SRCEXT      	= c
SRCDIR      	= src
OBJEXT      	= o
HEADEREXT		= hpp
HEADERDIR		= include
BUILDDIR    	= obj

LIBS			=
INCLUDES		= -I$(HEADERDIR)
LINK			=

SRC 			= $(shell find $(SRCDIR) -name '*.$(SRCEXT)')
HEADERS 		= $(shell find $(HEADERDIR) -name '*.$(HEADEREXT)')
OBJ				= $(foreach src,$(SRC),$(BUILDDIR)/$(notdir $(src:.$(SRCEXT)=.$(OBJEXT))))

VPATH 			= $(shell find $(SRCDIR) -type d | tr '\n' ':' | sed -E 's/(.*):/\1/')


# disable built-in rules
.SUFFIXES:

all: $(NAME)

$(NAME): $(BUILDDIR)/ $(OBJ) $(HEADERS)
	$(CC) $(CFLAGS) $(INCLUDES) $(OBJ) $(LIBS) -o $(NAME) $(LINK)

$(BUILDDIR)/%.$(OBJEXT): %.$(SRCEXT) $(HEADERS)
	$(CC) $(CFLAGS) $(INCLUDES) -c $< -o $(BUILDDIR)/$(notdir $@)

# Make sure that you run `make fclean` if you alredy compiled without debug
debug:
	DEBUG=1 $(MAKE)

clean:
ifneq ($(BUILDDIR),.)
	/bin/rm -rf $(BUILDDIR)/
endif

fclean: | clean
	/bin/rm -f $(NAME)

re: | fclean all

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)

silent:
	@$(MAKE) > /dev/null

.PHONY: all clean fclean re silent
