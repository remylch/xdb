Features : 
- [] Create HashIndex ( struct of filename, block, offset )
  ex : userID1234 → file7, block12, offset45
  index -> hash(index) -> (file, block, offset) -> data
- [] Create Data block (fixed size = 4kb) each block will be any data (json/string...) hashed too
