# import-layers-go

idea is to enforce abstraction levels / layers so that they dont mix up and you cna create coherent layered architecture



## linting & static checks


custom multichecker:

      cd cmd 
      go run . ../... &> ../reports/report.txt

or

      ./cmd/import-layers-go ./... &> reports/report.txt
