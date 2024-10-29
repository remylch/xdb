import ReactSelect from 'react-select';
import { classNames } from '../utils.ts';

type option = { label: string, value: string }


interface SelectProps {
  label?: string;
  placeholder?: string;
  options: option[];
  onChange: (value: any) => void;
  w?: "fit" | "full";
}

export const Select = ({ w = "full" ,...props }: SelectProps) => {
  return <div className={classNames('flex flex-col gap-2', w === "full" ? "w-full" : "w-fit")}>
    <label>{props.label}</label>
    <ReactSelect
      placeholder={props.placeholder}
      className={"w-full"}
      onChange={props.onChange}
      isSearchable={true}
      options={props.options}
      isMulti={false}
    />
  </div>;
};
